package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	svg "github.com/ajstarks/svgo"
)

var (
	width                  = 431
	xMargin                = 14
	xParagraph             = 35
	xFeeLine               = 55
	xMarginRightIndent     = 25
	xMarginRightIndentHard = 106
	xIndent                = 29
)

const (
	labelTitle                                 = "font-size:36pt;font-weight:900;font-family:'Roboto Flex';letter-spacing:0em"
	labelCompanyName                           = "font-size:18pt;letter-spacing:0em;font-weight:bold;font-family:'Roboto';text-anchor:left"
	labelPackageName                           = "font-size:14pt;letter-spacing:0em;font-weight:800;font-family:'Roboto Flex';text-anchor:left"
	labelGenericTextNormal                     = "font-size:12pt;letter-spacing:0em;font-family:Roboto;text-anchor:left"
	labelGenericTextNormalBold                 = "font-size:12pt;letter-spacing:0em;font-weight:900;font-family:Roboto;text-anchor:left"
	labelGenericTextNormalBoldAnchorEnd        = "font-size:12pt;letter-spacing:0em;font-weight:bold;font-family:'Roboto Flex';text-anchor:end"
	labelGenericTextNormalBoldAnchorStart      = "font-size:12pt;letter-spacing:0em;font-weight:bold;font-family:'Roboto Flex';text-anchor:start"
	labelGenericTextNormalHeavyBoldAnchorEnd   = "font-size:12pt;letter-spacing:0em;font-weight:900;font-family:Roboto;text-anchor:end"
	labelGenericTextNormalHeavyBoldAnchorStart = "font-size:12pt;letter-spacing:0em;font-weight:900;font-family:Roboto;text-anchor:start"
	labelGenericTextSmall                      = "font-size:10pt;letter-spacing:0em;font-family:Roboto;text-anchor:left"
	labelGenericTextSmallBoldAnchorEnd         = "font-size:10pt;letter-spacing:0em;font-weight:900;font-family:Roboto;text-anchor:end"
	labelMonthlyPrice                          = "font-size:18pt;letter-spacing:0em;font-weight:800;font-family:'Roboto Flex';text-anchor:left"
	labelMonthlyPriceValue                     = "font-size:18pt;letter-spacing:0em;font-weight:800;font-family:'Roboto Flex';text-anchor:end"
	labelSectionHeading                        = "font-size:14pt;letter-spacing:0em;font-weight:bold;font-family:'Roboto Flex';text-anchor:left"
	labelFccLink                               = "font-size:14pt;letter-spacing:0em;font-family:'Roboto Flex';text-anchor:end"
	labelUniquePlanId                          = "font-size:12pt;letter-spacing:0em;font-family:Roboto;text-anchor:left"
)

var templateStyle = `
    <style type="text/css">
       @import url('https://fonts.googleapis.com/css2?family=Roboto:wght@400;500;700;900');
       @import url('https://fonts.googleapis.com/css2?family=Roboto+Flex:opsz,wght@8..144,400;8..144,500;8..144,600;8..144,700;8..144,800;8..144,900;8..144,1000');                a:link,
       a:hover,
       a:active,
       a:visited {
           fill: #0000EE;
       }
	   a:hover {
	       text-decoration: underline;
       }
	</style>
`

func setTemplateStyles(canvas *svg.SVG) {
	canvas.Def()
	fmt.Fprintln(canvas.Writer, templateStyle)
	canvas.DefEnd()
}

type BroadbandConsumerLabel struct {
	yCounter int
}

func (b *BroadbandConsumerLabel) addY(offset int) int {
	b.yCounter += offset
	return b.yCounter
}

func (b *BroadbandConsumerLabel) getY() int {
	return b.yCounter
}

func (b *BroadbandConsumerLabel) labelTitle(canvas *svg.SVG, thisSectionYStart int) {
	canvas.Gstyle("")
	canvas.Text(xMargin, b.addY(55), "Broadband", labelTitle)
	canvas.Text(285, b.getY(), "Facts", labelTitle)
	canvas.Line(xMargin, b.addY(4), width-xMargin, b.getY(), "stroke:black;stroke-width:1")
	canvas.Gend()
}

func (b *BroadbandConsumerLabel) providerBlock(canvas *svg.SVG, thisSectionYStart int, template BroadbandData, fontList string) {
	canvas.Gstyle("")
	canvas.Text(xMargin, b.addY(25), template.CompanyName, labelCompanyName)
	canvas.Text(xMargin, b.addY(20), template.DataServiceName, labelPackageName)
	providerServiceType := ""
	if template.FixedOrMobile == "Fixed" {
		providerServiceType = "Fixed Broadband Consumer Disclosure"
	} else {
		providerServiceType = "Mobile Broadband Consumer Disclosure"
	}
	canvas.Text(xMargin, b.addY(21), providerServiceType, labelGenericTextNormal)
	canvas.Line(xMargin, b.addY(9), width-xMargin, b.getY(), "stroke:black;stroke-width:12")
	canvas.Gend()
}

func (b *BroadbandConsumerLabel) monthlyPrice(canvas *svg.SVG, thisSectionYStart int, template BroadbandData, fontList string) {
	canvas.Gstyle("")
	canvas.Text(xMargin, b.addY(30), "Monthly Price", labelMonthlyPrice)
	canvas.Text((width - xMargin), b.getY(), "$"+template.MonthlyPrice+"", labelMonthlyPriceValue)
	canvas.Line(xMargin, b.addY(8), width-xMargin, b.getY(), "stroke:black;stroke-width:3")
	canvas.Gend()
}

func (b *BroadbandConsumerLabel) monthlyDetails(canvas *svg.SVG, thisSectionYStart int, template BroadbandData, fontList string) {
	canvas.Gstyle("")
	// is introductory or not?
	if template.IntroductoryRate {
		canvas.Text(xMargin, b.addY(20), "This Monthly Price is an introductory rate.", labelGenericTextNormal)
		lineY := b.addY(14)
		canvas.Text(xIndent, lineY, "Introductory Period", labelGenericTextNormal)
		canvas.Text((width - xMargin), lineY, template.IntroductoryPeriodInMonths+" months", labelGenericTextNormalBoldAnchorEnd)
		lineY = b.addY(17)
		canvas.Text(xIndent, lineY, "Price after introductory period", labelGenericTextNormal)
		canvas.Text((width - xMargin), lineY, "$"+template.DataServicePrice, labelGenericTextNormalBoldAnchorEnd)
		contractDuration, err := strconv.Atoi(template.ContractDuration)
		if err != nil {
			log.Fatalln("error:", err)
			return
		}
		aOrAn := ""
		if contractDuration == 8 || contractDuration == 11 || contractDuration == 18 {
			aOrAn = "an"
		} else {
			aOrAn = "a"
		}
		contractTerms := "This Monthly Price requires " + aOrAn + " " + template.ContractDuration + " month"
		canvas.Textspan(xMargin, b.addY(17), contractTerms, labelGenericTextNormal)
		canvas.Link(template.ContractURL, "contract")
		// css for blue a links
		canvas.Span("contract", "fill:blue")
		canvas.LinkEnd()
		canvas.TextEnd()
	} else {
		canvas.Text(xMargin, b.addY(17), "This Monthly Price is not an introductory rate.", labelGenericTextNormal)
		canvas.Text(xMargin, b.addY(17), "This Monthly Price does not require a contract.", labelGenericTextNormal)
	}
	lineY := b.addY(12)
	canvas.Line(xMargin, lineY, width-xMargin, lineY, "stroke:black;stroke-width:1")
	canvas.Gend()
}

// TODO limit to 37 characters or less..
func (b *BroadbandConsumerLabel) additionalChargesAndTerms(canvas *svg.SVG, thisSectionYStart int, template BroadbandData, fontList string) {
	canvas.Gstyle("")
	canvas.Text(xMargin, b.addY(23), "Additional Charges & Terms", labelSectionHeading)

	canvas.Text(xParagraph, b.addY(17), "Provider Monthly Fees", labelGenericTextNormal)
	if len(template.ExtraMonthlyFields) > 0 {
		for _, charge := range template.ExtraMonthlyFields {
			charge.ChargeValue = strings.TrimPrefix(charge.ChargeValue, "$")
			canvas.Text(xFeeLine, b.addY(17), charge.ChargeName, labelGenericTextNormal)
			canvas.Text((width - xMarginRightIndent), b.getY(), "$"+charge.ChargeValue, labelGenericTextNormalBoldAnchorEnd)
		}
	} else {
		canvas.Text(xFeeLine, b.addY(17), "No additional monthly fees", labelGenericTextNormal)
	}

	canvas.Text(xParagraph, b.addY(35), "One-time Fees at the Time of Purchase", labelGenericTextNormal)
	if len(template.ExtraOneTimeFields) > 0 {
		for _, charge := range template.ExtraOneTimeFields {
			charge.ChargeValue = strings.TrimPrefix(charge.ChargeValue, "$")
			canvas.Text(xFeeLine, b.addY(17), charge.ChargeName, labelGenericTextNormal)
			canvas.Text((width - xMarginRightIndent), b.getY(), "$"+charge.ChargeValue, labelGenericTextNormalBoldAnchorEnd)
		}
	} else {
		canvas.Text(xFeeLine, b.addY(17), "No additional one-time fees at time of purchase", labelGenericTextNormal)

	}

	if template.EarlyTerminationFee != "" {
		template.EarlyTerminationFee = strings.TrimPrefix(template.EarlyTerminationFee, "$")
		canvas.Text(xParagraph, b.addY(35), "Early Termination Fee", labelGenericTextNormal)
		canvas.Text((width - xMarginRightIndent), b.getY(), "$"+template.EarlyTerminationFee, labelGenericTextNormalBoldAnchorEnd)
	} else {
		canvas.Text(xParagraph, b.addY(35), "Early Termination Fee", labelGenericTextNormal)
		canvas.Text((width - xMarginRightIndent), b.getY(), "None", labelGenericTextNormalHeavyBoldAnchorEnd)
	}
	canvas.Text(xParagraph, b.addY(35), "Government Taxes", labelGenericTextNormal)
	canvas.Text((width - xMargin), b.getY(), "Varies by Location", labelGenericTextNormalHeavyBoldAnchorEnd)
	canvas.Line(xMargin, b.addY(10), width-xMargin, b.getY(), "stroke:black;stroke-width:3")
	canvas.Gend()
}

func (b *BroadbandConsumerLabel) discountsAndBundles(canvas *svg.SVG, thisSectionYStart int, template BroadbandData, fontList string) {
	canvas.Gstyle("")
	canvas.Text(xMargin, b.addY(23), "Discounts & Bundles", labelSectionHeading)
	canvas.Textspan(xParagraph, b.addY(17), "", labelGenericTextNormal)
	canvas.Link(template.DiscountsAndBundlesURL, "Click Here")
	canvas.Span("Click here", "fill:blue")
	canvas.LinkEnd()
	canvas.TextEnd()
	canvas.Text(110, b.getY(), "for available billing discounts and pricing", labelGenericTextNormal)
	canvas.Text(xParagraph, b.addY(17), "options for broadband service bundled with other", labelGenericTextNormal)
	canvas.Text(xParagraph, b.addY(17), "services like video, phone, and wireless service", labelGenericTextNormal)
	canvas.Text(xParagraph, b.addY(17), "and use of your own equipment like modems and", labelGenericTextNormal)
	canvas.Text(xParagraph, b.addY(17), "routers.", labelGenericTextNormal)
	canvas.Line(xMargin, b.addY(10), width-xMargin, b.getY(), "stroke:black;stroke-width:1")
	canvas.Gend()
}

func (b *BroadbandConsumerLabel) participatesInACP(canvas *svg.SVG, thisSectionYStart int, template BroadbandData, fontList string) {
	canvas.Gstyle("")
	canvas.Text(xMargin, b.addY(23), "Affordable Connectivity Program (ACP)", labelSectionHeading)
	canvas.Text(xParagraph, b.addY(17), "The ACP is a government program to help lower the", labelGenericTextNormal)
	canvas.Text(xParagraph, b.addY(17), "monthly cost of internet service. To learn more", labelGenericTextNormal)
	canvas.Text(xParagraph, b.addY(17), "about the ACP, including to find out whether you", labelGenericTextNormal)
	canvas.Text(xParagraph, b.addY(17), "qualify, visit:", labelGenericTextNormal)
	canvas.Textspan(130, b.getY(), "", labelGenericTextNormal)
	canvas.Link("https://affordableconnectivity.gov/", "affordableconnectivity.gov")
	canvas.Span("affordableconnectivity.gov", "fill:blue")
	canvas.LinkEnd()
	canvas.TextEnd()

	template.AcpEnabled = strings.ToUpper(template.AcpEnabled)
	if template.AcpEnabled == "YES" || template.AcpEnabled == "1" || template.AcpEnabled == "TRUE" {
		canvas.Text(56, b.addY(17), "Participates in the ACP", labelGenericTextNormalBold)
		canvas.Text((width - xMarginRightIndentHard), b.getY(), "Yes", labelGenericTextNormalHeavyBoldAnchorStart)
	} else {
		canvas.Text(56, b.addY(17), "Participates in the ACP", labelGenericTextNormalBold)
		canvas.Text((width - xMarginRightIndentHard), b.getY(), "No", labelGenericTextNormalHeavyBoldAnchorStart)
	}
	canvas.Line(xMargin, b.addY(10), width-xMargin, b.getY(), "stroke:black;stroke-width:3")
	canvas.Gend()
}

func (b *BroadbandConsumerLabel) planSpeeds(canvas *svg.SVG, thisSectionYStart int, template BroadbandData, fontList string) {
	canvas.Gstyle("")
	canvas.Text(xMargin, b.addY(23), "Speeds Provided with Plan", labelSectionHeading)
	canvas.Text(xIndent, b.addY(17), "Typical Download Speed", labelGenericTextNormal)
	canvas.Text((width - xMarginRightIndentHard), b.getY(), template.CalculatedDLSpeedInMbps+" Mbps", labelGenericTextNormalHeavyBoldAnchorStart)
	canvas.Text(xIndent, b.addY(17), "Typical Upload Speed", labelGenericTextNormal)
	canvas.Text((width - xMarginRightIndentHard), b.getY(), template.CalculatedULSpeedInMbps+" Mbps", labelGenericTextNormalHeavyBoldAnchorStart)
	canvas.Text(xIndent, b.addY(17), "Typical Latency", labelGenericTextNormal)
	canvas.Text((width - xMarginRightIndentHard), b.getY(), template.LatencyInMs+" ms", labelGenericTextNormalHeavyBoldAnchorStart)
	canvas.Line(xMargin, b.addY(10), width-xMargin, b.getY(), "stroke:black;stroke-width:1")
	canvas.Gend()

}

func (b *BroadbandConsumerLabel) dataIncluded(canvas *svg.SVG, thisSectionYStart int, template BroadbandData, fontList string) {
	canvas.Gstyle("")
	if template.DataIncludedInMonthlyPriceGB != "" {
		canvas.Text(xMargin, b.addY(23), "Data Included with Monthly Price", labelSectionHeading)
		canvas.Text((width - xMarginRightIndentHard), b.getY(), template.DataIncludedInMonthlyPriceGB+" GB", labelGenericTextNormalHeavyBoldAnchorStart)
		canvas.Text(xIndent, b.addY(17), "Charges for Additional Data Usage", labelGenericTextNormal)
		if template.OverageFee == "" {
			canvas.Text((width - xMarginRightIndentHard), b.getY(), "None", labelGenericTextNormalHeavyBoldAnchorStart)
		} else {
			template.OverageFee = strings.TrimPrefix(template.OverageFee, "$")
			canvas.Text((width - xMarginRightIndentHard), b.getY(), "$"+template.OverageFee+"/"+template.OverageDataAmount+"GB", labelGenericTextNormalHeavyBoldAnchorStart)
		}
	} else {
		canvas.Text(xMargin, b.addY(23), "Data Included with Monthly Price", labelSectionHeading)
		canvas.Text((width - xMarginRightIndentHard), b.getY(), "Unlimited", labelGenericTextNormalHeavyBoldAnchorStart)
		canvas.Text(xIndent, b.addY(17), "Charges for Additional Data Usage", labelGenericTextNormal)
		canvas.Text((width - xMarginRightIndentHard), b.getY(), "None", labelGenericTextNormalHeavyBoldAnchorStart)
	}
	canvas.Line(xMargin, b.addY(10), width-xMargin, b.getY(), "stroke:black;stroke-width:3")
	canvas.Gend()
}

func (b *BroadbandConsumerLabel) policies(canvas *svg.SVG, thisSectionYStart int, template BroadbandData, fontList string) {
	canvas.Gstyle("")
	canvas.Text(xMargin, b.addY(23), "Network Management", labelSectionHeading)
	canvas.Textspan(width-xMargin, b.getY(), "", labelGenericTextNormalBoldAnchorEnd)
	canvas.Link(template.NetworkManagementURL, template.NetworkManagementURL)
	canvas.Span("Read our Policy", "fill:blue")
	canvas.LinkEnd()
	canvas.TextEnd()

	canvas.Text(xMargin, b.addY(17), "Privacy", labelSectionHeading)
	canvas.Textspan(width-xMargin, b.getY(), "", labelGenericTextNormalBoldAnchorEnd)
	canvas.Link(template.PrivacyPolicyURL, template.PrivacyPolicyURL)
	canvas.Span("Read our Policy", "fill:blue")
	canvas.LinkEnd()
	canvas.TextEnd()
	canvas.Line(xMargin, b.addY(15), width-xMargin, b.getY(), "stroke:black;stroke-width:12")
	canvas.Gend()
}

func (b *BroadbandConsumerLabel) customerSupport(canvas *svg.SVG, thisSectionYStart int, template BroadbandData, fontList string) {
	canvas.Gstyle("")
	canvas.Text(xMargin, b.addY(25), "Customer Support", labelSectionHeading)

	canvas.Text(xIndent, b.addY(17), "Contact Us:", labelGenericTextNormal)
	canvas.Textspan(width-310, b.getY(), "", labelGenericTextNormal)
	canvas.Link(template.CustomerSupportURL, template.CustomerSupportURL)
	canvas.Span("Contact Us", "fill:blue")
	canvas.LinkEnd()
	canvas.TextEnd()
	canvas.Text(width-220, b.getY(), "/ "+template.CustomerSupportPhone, labelGenericTextNormal)
	canvas.Line(xMargin, b.addY(15), width-xMargin, b.getY(), "stroke:black;stroke-width:6")
	canvas.Gend()
}

func (b *BroadbandConsumerLabel) fccLabelTerms(canvas *svg.SVG, thisSectionYStart int, template BroadbandData, fontList string) {
	canvas.Gstyle("")
	canvas.Text(xMargin, b.addY(17), "Learn more about the terms used on this label by", labelGenericTextNormal)
	canvas.Text(xMargin, b.addY(17), "visiting the Federal Communications Commission's", labelGenericTextNormal)
	canvas.Text(xMargin, b.addY(17), "Consumer Resource Center.", labelGenericTextNormal)

	canvas.Textspan(width-xMargin, b.addY(17), "", labelFccLink)
	canvas.Link("https://fcc.gov/consumer", "https://fcc.gov/consumer")
	canvas.Span("fcc.gov/consumer", "fill:blue")
	canvas.LinkEnd()
	canvas.TextEnd()
	canvas.Gend()
}

func (b *BroadbandConsumerLabel) uniquePlanIdentifier(canvas *svg.SVG, thisSectionYStart int, template BroadbandData, fontList string) {
	canvas.Gstyle("")
	FixedMobile := ""
	if template.FixedOrMobile == "Fixed" {
		FixedMobile = "F"
	} else {
		FixedMobile = "M"
	}
	serviceId := strings.Repeat("0", 15-len(template.DataServiceID)) + template.DataServiceID
	uniquePlanIdentifier := FixedMobile + template.FccID + serviceId

	canvas.Text(xMargin, b.addY(17), uniquePlanIdentifier, labelUniquePlanId)
	canvas.Line(xMargin, b.addY(15), width-xMargin, b.getY(), "stroke:white;stroke-width:6")
	canvas.Gend()
}

var svgStartTag = `
<!-- coded by andy, katherine and gene @ sonar.software -->
<!-- https://www.sonar.software -->

<svg
     id="bcd"
     viewBox="{{ .TemplateViewBox }}"
     xmlns="http://www.w3.org/2000/svg"
     xmlns:xlink="http://www.w3.org/1999/xlink">`

func generateLabels(templateData []BroadbandData) error {

	for templateNumber, template := range templateData {
		templateFile, err := os.Create(fmt.Sprintf("%v/label_%d.svg", outputDirectory, templateNumber))
		if err != nil {
			log.Fatalln("error:", err)
			return err
		}
		defer templateFile.Close()

		templateWriter := &TemplateWriter{}

		canvas := svg.New(templateWriter)
		fmt.Fprintln(templateWriter, svgStartTag)

		setTemplateStyles(canvas)
		fontList := "Roboto"

		var label BroadbandConsumerLabel

		canvas.Gid("content-group")
		fmt.Fprintln(templateWriter, `<rect width="431" height="942" style="fill:white" />`)
		fmt.Fprintln(templateWriter, `<rect x="4.5" y="7.5" width="419" height="{{ .CalcYRectHeight }}" style="fill:none;stroke:black;stroke-width:3" />`)

		label.labelTitle(canvas, 0)
		label.providerBlock(canvas, 0, template, fontList)
		label.monthlyPrice(canvas, 0, template, fontList)
		label.monthlyDetails(canvas, 0, template, fontList)
		label.additionalChargesAndTerms(canvas, 0, template, fontList)
		label.discountsAndBundles(canvas, 0, template, fontList)
		label.participatesInACP(canvas, 0, template, fontList)
		label.planSpeeds(canvas, 0, template, fontList)
		label.dataIncluded(canvas, 0, template, fontList)
		label.policies(canvas, 0, template, fontList)
		label.customerSupport(canvas, 0, template, fontList)
		label.fccLabelTerms(canvas, 0, template, fontList)
		label.uniquePlanIdentifier(canvas, 0, template, fontList)
		canvas.Gend()
		canvas.End()

		// bake the viewbox

		templateWriter.ApplyDynamicCalculations(label.getY())

		// write the contents of templateWriter to the templateFile
		for _, v := range templateWriter.bcdTemplate {
			fmt.Fprintln(templateFile, v)
		}

	}
	return nil
}

// TemplateWriter is a custom io.Writer that appends data to a []string.
type TemplateWriter struct {
	bcdTemplate []string
}

func (w *TemplateWriter) Write(p []byte) (n int, err error) {
	w.bcdTemplate = append(w.bcdTemplate, string(p))
	return len(p), nil
}

func (w *TemplateWriter) ApplyDynamicCalculations(y int) {
	for i, v := range w.bcdTemplate {
		if strings.Contains(v, "{{ .TemplateViewBox }}") {
			viewBox := "0 0 431 " + strconv.Itoa(y)
			w.bcdTemplate[i] = strings.ReplaceAll(v, "{{ .TemplateViewBox }}", viewBox)
		}
		if strings.Contains(v, "{{ .CalcYRectHeight }}") {
			rectHeight := strconv.Itoa(y - 13)
			w.bcdTemplate[i] = strings.ReplaceAll(v, "{{ .CalcYRectHeight }}", rectHeight)
		}
	}

}
