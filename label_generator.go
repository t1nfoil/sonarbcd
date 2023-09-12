package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	svg "github.com/ajstarks/svgo"
)

var (
	width = 605
)

const (
	labelTitle                               = "font-size:48pt;font-weight:945;font-family:'Roboto Flex';text-anchor:middle"
	labelCompanyName                         = "font-size:18pt;font-weight:bold;font-family:'Roboto';text-anchor:left"
	labelPackageName                         = "font-size:14pt;font-weight:800;font-family:'Roboto Flex';text-anchor:left"
	labelGenericTextNormal                   = "font-size:14pt;font-family:Roboto;text-anchor:left"
	labelGenericTextNormalBold               = "font-size:14pt;font-weight:900;font-family:Roboto;text-anchor:left"
	labelGenericTextNormalBoldAnchorEnd      = "font-size:14pt;font-weight:bold;font-family:'Roboto Flex';text-anchor:end"
	labelGenericTextNormalHeavyBoldAnchorEnd = "font-size:14pt;font-weight:900;font-family:Roboto;text-anchor:end"
	labelGenericTextSmall                    = "font-size:12pt;font-family:Roboto;text-anchor:left"
	labelGenericTextSmallBoldAnchorEnd       = "font-size:12pt;font-weight:900;font-family:Roboto;text-anchor:end"
	labelMonthlyPrice                        = "font-size:18pt;font-weight:800;font-family:'Roboto Flex';text-anchor:left"
	labelMonthlyPriceValue                   = "font-size:18pt;font-weight:800;font-family:'Roboto Flex';text-anchor:end"
	labelSectionHeading                      = "font-size:14pt;font-weight:bold;font-family:'Roboto Flex';text-anchor:left"
	labelFccLink                             = "font-size:14pt;font-family:'Roboto Flex';text-anchor:end"
	labelUniquePlanId                        = "font-size:12pt;font-family:Roboto;text-anchor:left"
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

func (b *BroadbandConsumerLabel) labelTitle(canvas *svg.SVG, thisSectionYStart int, title string, fontList string) {
	canvas.Gstyle("")
	canvas.Text(width/2, b.addY(65), title, labelTitle)
	canvas.Line(25, b.addY(10), width-25, b.getY(), "stroke:black;stroke-width:1")
	canvas.Gend()
}

func (b *BroadbandConsumerLabel) providerBlock(canvas *svg.SVG, thisSectionYStart int, template BroadbandData, fontList string) {
	canvas.Gstyle("")
	canvas.Text(25, b.addY(30), template.CompanyName, labelCompanyName)
	canvas.Text(25, b.addY(30), template.DataServiceName, labelPackageName)
	providerServiceType := ""
	if template.FixedOrMobile == "Fixed" {
		providerServiceType = "Fixed Broadband Consumer Disclosure"
	} else {
		providerServiceType = "Mobile Broadband Consumer Disclosure"
	}
	canvas.Text(25, b.addY(30), providerServiceType, labelGenericTextNormal)
	canvas.Line(25, b.addY(10), width-25, b.getY(), "stroke:black;stroke-width:12")
	canvas.Gend()
}

func (b *BroadbandConsumerLabel) monthlyPrice(canvas *svg.SVG, thisSectionYStart int, template BroadbandData, fontList string) {
	canvas.Gstyle("")
	canvas.Text(25, b.addY(30), "Monthly Price", labelMonthlyPrice)
	canvas.Text((width - 25), b.getY(), "$"+template.MonthlyPrice+"", labelMonthlyPriceValue)
	canvas.Line(25, b.addY(10), width-25, b.getY(), "stroke:black;stroke-width:3")
	canvas.Gend()
}

func (b *BroadbandConsumerLabel) monthlyDetails(canvas *svg.SVG, thisSectionYStart int, template BroadbandData, fontList string) {
	canvas.Gstyle("")
	// is introductory or not?
	if template.IntroductoryRate {
		canvas.Text(25, b.addY(25), "This Monthly Price is an introductory rate.", labelGenericTextNormal)
		lineY := b.addY(25)
		canvas.Text(40, lineY, "Introductory Period", labelGenericTextSmall)
		canvas.Text((width - 25), lineY, template.IntroductoryPeriodInMonths+" months", labelGenericTextSmallBoldAnchorEnd)
		lineY = b.addY(25)
		canvas.Text(40, lineY, "Price after introductory period", labelGenericTextSmall)
		canvas.Text((width - 25), lineY, "$"+template.DataServicePrice, labelGenericTextSmallBoldAnchorEnd)
		contractDuration, err := strconv.Atoi(template.ContractDuration)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		aOrAn := ""
		if contractDuration == 8 || contractDuration == 11 || contractDuration == 18 {
			aOrAn = "an"
		} else {
			aOrAn = "a"
		}
		contractTerms := "This Monthly Price requires " + aOrAn + " " + template.ContractDuration + " month"
		canvas.Textspan(25, b.addY(25), contractTerms, labelGenericTextNormal)
		canvas.Link(template.ContractURL, "contract")
		// css for blue a links
		canvas.Span("contract", "fill:blue")
		canvas.LinkEnd()
		canvas.TextEnd()
	} else {
		canvas.Text(25, b.addY(25), "This Monthly Price is not an introductory rate.", labelGenericTextNormal)
		canvas.Text(25, b.addY(25), "This Monthly Price does not require a contract.", labelGenericTextNormal)
	}
	lineY := b.addY(10)
	canvas.Line(25, lineY, width-25, lineY, "stroke:black;stroke-width:1")
	canvas.Gend()
}

// TODO limit to 37 characters or less..
func (b *BroadbandConsumerLabel) additionalChargesAndTerms(canvas *svg.SVG, thisSectionYStart int, template BroadbandData, fontList string) {
	canvas.Gstyle("")
	canvas.Text(25, b.addY(25), "Additional Charges & Terms", labelSectionHeading)

	canvas.Text(25, b.addY(35), "Provider Monthly Fees", labelGenericTextNormal)
	if len(template.ExtraMonthlyFields) > 0 {
		for _, charge := range template.ExtraMonthlyFields {
			charge.ChargeValue = strings.TrimPrefix(charge.ChargeValue, "$")
			canvas.Text(55, b.addY(25), charge.ChargeName, labelGenericTextSmall)
			canvas.Text((width - 25), b.getY(), "$"+charge.ChargeValue, labelGenericTextSmallBoldAnchorEnd)
		}
	} else {
		canvas.Text(55, b.addY(25), "No additional monthly fees", labelGenericTextNormal)
	}

	canvas.Text(25, b.addY(35), "One-time Fees at the Time of Purchase", labelGenericTextNormal)
	if len(template.ExtraOneTimeFields) > 0 {
		for _, charge := range template.ExtraOneTimeFields {
			charge.ChargeValue = strings.TrimPrefix(charge.ChargeValue, "$")
			canvas.Text(55, b.addY(25), charge.ChargeName, labelGenericTextSmall)
			canvas.Text((width - 25), b.getY(), "$"+charge.ChargeValue, labelGenericTextSmallBoldAnchorEnd)
		}
	} else {
		canvas.Text(55, b.addY(25), "No additional one-time fees at time of purchase", labelGenericTextNormal)

	}

	if template.EarlyTerminationFee != "" {
		template.EarlyTerminationFee = strings.TrimPrefix(template.EarlyTerminationFee, "$")
		canvas.Text(25, b.addY(35), "Early Termination Fee", labelGenericTextNormal)
		canvas.Text((width - 25), b.getY(), "$"+template.EarlyTerminationFee, labelGenericTextNormalHeavyBoldAnchorEnd)
	} else {
		canvas.Text(25, b.addY(35), "Early Termination Fee", labelGenericTextNormal)
		canvas.Text((width - 25), b.getY(), "None", labelGenericTextNormalHeavyBoldAnchorEnd)
	}
	canvas.Text(25, b.addY(35), "Government Taxes", labelGenericTextNormal)
	canvas.Text((width - 25), b.getY(), "Varies by Location", labelGenericTextNormalHeavyBoldAnchorEnd)
	canvas.Line(25, b.addY(10), width-25, b.getY(), "stroke:black;stroke-width:3")
	canvas.Gend()
}

func (b *BroadbandConsumerLabel) discountsAndBundles(canvas *svg.SVG, thisSectionYStart int, template BroadbandData, fontList string) {
	canvas.Gstyle("")
	canvas.Text(25, b.addY(25), "Discounts & Bundles", labelSectionHeading)
	canvas.Textspan(40, b.addY(25), "", labelGenericTextNormal)
	canvas.Link(template.DiscountsAndBundlesURL, "Click Here")
	canvas.Span("Click here", "fill:blue")
	canvas.LinkEnd()
	canvas.TextEnd()
	canvas.Text(130, b.getY(), " for available billing discounts and pricing options", labelGenericTextNormal)
	canvas.Text(40, b.addY(25), "for broadband service bundled with other services like video,", labelGenericTextNormal)
	canvas.Text(40, b.addY(25), "phone, and wireless service, and use of your own equipment", labelGenericTextNormal)
	canvas.Text(40, b.addY(25), "like modems and routers.", labelGenericTextNormal)
	canvas.Line(25, b.addY(10), width-25, b.getY(), "stroke:black;stroke-width:3")
	canvas.Gend()
}

func (b *BroadbandConsumerLabel) participatesInACP(canvas *svg.SVG, thisSectionYStart int, template BroadbandData, fontList string) {
	canvas.Gstyle("")
	canvas.Text(25, b.addY(25), "Affordable Connectivity Program (ACP)", labelSectionHeading)
	canvas.Text(40, b.addY(25), "The ACP is a government program to help lower the monthly", labelGenericTextNormal)
	canvas.Text(40, b.addY(25), "cost of internet service. To learn more about the ACP, including", labelGenericTextNormal)
	canvas.Text(40, b.addY(25), "to find out whether you qualify, visit:", labelGenericTextNormal)
	canvas.Textspan(340, b.getY(), "", labelGenericTextNormal)
	canvas.Link("https://affordableconnectivity.gov/", "affordableconnectivity.gov")
	canvas.Span("affordableconnectivity.gov", "fill:blue")
	canvas.LinkEnd()
	canvas.TextEnd()

	template.AcpEnabled = strings.ToUpper(template.AcpEnabled)
	if template.AcpEnabled == "YES" || template.AcpEnabled == "1" || template.AcpEnabled == "TRUE" {
		canvas.Text(55, b.addY(25), "Participates in the ACP", labelGenericTextNormalBold)
		canvas.Text((width - 25), b.getY(), "Yes", labelGenericTextNormalHeavyBoldAnchorEnd)
	} else {
		canvas.Text(55, b.addY(25), "Participates in the ACP", labelGenericTextNormalBold)
		canvas.Text((width - 25), b.getY(), "No", labelGenericTextNormalHeavyBoldAnchorEnd)
	}
	canvas.Line(25, b.addY(10), width-25, b.getY(), "stroke:black;stroke-width:3")
	canvas.Gend()
}

func (b *BroadbandConsumerLabel) planSpeeds(canvas *svg.SVG, thisSectionYStart int, template BroadbandData, fontList string) {
	canvas.Gstyle("")
	canvas.Text(25, b.addY(25), "Plan Speeds", labelSectionHeading)
	canvas.Text(40, b.addY(25), "Typical Download Speed", labelGenericTextNormal)
	canvas.Text((width - 25), b.getY(), template.CalculatedDLSpeedInMbps+" Mbps", labelGenericTextNormalHeavyBoldAnchorEnd)
	canvas.Text(40, b.addY(25), "Typical Upload Speed", labelGenericTextNormal)
	canvas.Text((width - 25), b.getY(), template.CalculatedULSpeedInMbps+" Mbps", labelGenericTextNormalHeavyBoldAnchorEnd)
	canvas.Text(40, b.addY(25), "Typical Latency", labelGenericTextNormal)
	canvas.Text((width - 25), b.getY(), template.LatencyInMs+" ms", labelGenericTextNormalHeavyBoldAnchorEnd)
	canvas.Line(25, b.addY(10), width-25, b.getY(), "stroke:black;stroke-width:3")
	canvas.Gend()

}

func (b *BroadbandConsumerLabel) dataIncluded(canvas *svg.SVG, thisSectionYStart int, template BroadbandData, fontList string) {
	canvas.Gstyle("")
	if template.DataIncludedInMonthlyPriceGB != "" {
		canvas.Text(25, b.addY(25), "Data Included with Monthly Price", labelSectionHeading)
		canvas.Text((width - 25), b.getY(), template.DataIncludedInMonthlyPriceGB+" GB", labelGenericTextNormalBoldAnchorEnd)
		canvas.Text(40, b.addY(25), "Charges for Additional Data Usage", labelGenericTextNormal)
		template.OverageFee = strings.TrimPrefix(template.OverageFee, "$")
		canvas.Text((width - 25), b.getY(), "$"+template.OverageFee+"/"+template.OverageDataAmount+"GB", labelGenericTextNormalHeavyBoldAnchorEnd)
	} else {
		canvas.Text(25, b.addY(25), "Data Included with Monthly Price", labelSectionHeading)
		canvas.Text((width - 25), b.getY(), "Unlimited", labelGenericTextNormalBoldAnchorEnd)
		canvas.Text(40, b.addY(25), "Charges for Additional Data Usage", labelGenericTextNormal)
		canvas.Text((width - 25), b.getY(), "None", labelGenericTextNormalHeavyBoldAnchorEnd)
	}
	canvas.Line(25, b.addY(10), width-25, b.getY(), "stroke:black;stroke-width:3")
	canvas.Gend()
}

func (b *BroadbandConsumerLabel) policies(canvas *svg.SVG, thisSectionYStart int, template BroadbandData, fontList string) {
	canvas.Gstyle("")
	canvas.Text(25, b.addY(25), "Network Management", labelSectionHeading)
	canvas.Textspan(width-25, b.getY(), "", labelGenericTextNormalBoldAnchorEnd)
	canvas.Link(template.NetworkManagementURL, template.NetworkManagementURL)
	canvas.Span("Read our Policy", "fill:blue")
	canvas.LinkEnd()
	canvas.TextEnd()

	canvas.Text(25, b.addY(25), "Privacy", labelSectionHeading)
	canvas.Textspan(width-25, b.getY(), "", labelGenericTextNormalBoldAnchorEnd)
	canvas.Link(template.PrivacyPolicyURL, template.PrivacyPolicyURL)
	canvas.Span("Read our Policy", "fill:blue")
	canvas.LinkEnd()
	canvas.TextEnd()
	canvas.Line(25, b.addY(15), width-25, b.getY(), "stroke:black;stroke-width:12")
	canvas.Gend()
}

func (b *BroadbandConsumerLabel) customerSupport(canvas *svg.SVG, thisSectionYStart int, template BroadbandData, fontList string) {
	canvas.Gstyle("")
	canvas.Text(25, b.addY(25), "Customer Support", labelSectionHeading)

	canvas.Text(40, b.addY(25), "Contact Us:", labelGenericTextNormal)
	canvas.Textspan(width-390, b.getY(), "", labelGenericTextNormal)
	canvas.Link(template.CustomerSupportURL, template.CustomerSupportURL)
	canvas.Span("Contact Us", "fill:blue")
	canvas.LinkEnd()
	canvas.TextEnd()
	canvas.Text(width-240, b.getY(), "/ "+template.CustomerSupportPhone, labelGenericTextNormal)
	canvas.Line(25, b.addY(15), width-25, b.getY(), "stroke:black;stroke-width:6")
	canvas.Gend()
}

func (b *BroadbandConsumerLabel) fccLabelTerms(canvas *svg.SVG, thisSectionYStart int, template BroadbandData, fontList string) {
	canvas.Gstyle("")
	canvas.Text(25, b.addY(25), "Learn more about the terms used on this label by visiting the", labelGenericTextNormal)
	canvas.Text(25, b.addY(25), "Federal Communications Commission's Consumer Resource", labelGenericTextNormal)
	canvas.Text(25, b.addY(25), " Center.", labelGenericTextNormal)

	canvas.Textspan(width-25, b.addY(25), "", labelFccLink)
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

	canvas.Text(25, b.addY(25), uniquePlanIdentifier, labelUniquePlanId)
	canvas.Gend()
}

var resizeJsScript = `
	window.addEventListener('load', function() {
	        var contentGroup = document.getElementById('content-group');
	        var bbox = contentGroup.getBBox();
	        var svg = document.getElementById('bcd');
	        svg.setAttribute('height', bbox.height);
	});
`

func generateLabels(templateData []BroadbandData) error {

	// loop through templateData and open an svg file for each as an io.Writer
	for templateNumber, template := range templateData {
		templateFile, err := os.Create(fmt.Sprintf("%v/label_%d.svg", outputDirectory, templateNumber))
		if err != nil {
			fmt.Println("error:", err)
			return err
		}
		defer templateFile.Close()

		templateWriter := io.Writer(templateFile)

		canvas := svg.New(templateWriter)
		canvas.Startpercent(100, 100, "id=\"bcd\"")
		canvas.Gid("content-group")
		setTemplateStyles(canvas)
		fontList := "Roboto"
		canvas.Script("application/javascript", resizeJsScript)

		var label BroadbandConsumerLabel

		fmt.Fprintln(templateWriter, `<rect x="0" y="0" width="605" height="100%" style="fill:white" />`)
		fmt.Fprintln(templateWriter, `<rect x="5" y="5" width="595" height="99.5%" style="fill:none;stroke:black;stroke-width:3" />`)

		label.labelTitle(canvas, 0, "Broadband Facts", fontList)
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
	}
	return nil
}
